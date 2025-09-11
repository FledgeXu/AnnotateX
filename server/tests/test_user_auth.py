from datetime import datetime
from typing import Optional

import pytest
from returns.result import Failure, Success
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from annotatex.models.user import AuthIdentities, User
from annotatex.repositories.user_auth import UserAuth


@pytest.fixture()
def user_auth_repo(session: AsyncSession) -> UserAuth:
    """Provide a UserAuth repository bound to the test session."""
    return UserAuth(session)


async def _create_user(session: AsyncSession, username: str = "alice") -> User:
    """
    Create a minimal User.
    If your User model requires more fields, extend this helper accordingly.
    """
    user = User(username=username)
    async with session.begin():
        session.add(user)
        await session.flush()
    return user


async def _create_identity(
    repo: UserAuth,
    *,
    user_id,
    provider: str = "password",
    subject: str = "subject",
    email: Optional[str] = None,
    email_verified: bool = False,
    hash_method: Optional[str] = None,
    password_hash: Optional[str] = None,
) -> AuthIdentities:
    """
    Create an AuthIdentities row through the repository and unwrap the result.
    """
    res = await repo.add_auth_identity(
        user_id=user_id,
        provider=provider,
        subject=subject,
        email=email,
        email_verified=email_verified,
        hash_method=hash_method,
        password_hash=password_hash,
    )
    assert isinstance(res, Success)
    return res.unwrap()


class TestUserAuth:
    async def test_add_auth_identity_success(
        self, session: AsyncSession, user_auth_repo: UserAuth
    ):
        user = await _create_user(session)
        identity = await _create_identity(
            user_auth_repo,
            user_id=user.id,
            provider="password",
            subject="alice",
            email="alice@example.com",
            email_verified=False,
            hash_method="argon2",
            password_hash="argon2$dummy",
        )
        assert identity.id is not None

        # Verify it exists in DB
        row = (
            await session.execute(
                select(AuthIdentities).where(AuthIdentities.id == identity.id)
            )
        ).scalar_one_or_none()
        assert row is not None
        assert row.user_id == user.id
        assert row.provider == "password"
        assert row.subject == "alice"

    async def test_add_auth_identity_unique_conflict(
        self, session: AsyncSession, user_auth_repo: UserAuth
    ):
        """
        Assumes a unique constraint on (provider, subject) or equivalent.
        The second insertion with the same combination should fail.
        """
        user = await _create_user(session)

        first = await user_auth_repo.add_auth_identity(
            user_id=user.id, provider="password", subject="dup"
        )
        assert isinstance(first, Success)

        second = await user_auth_repo.add_auth_identity(
            user_id=user.id, provider="password", subject="dup"
        )
        assert isinstance(second, Failure)

    async def test_get_user_by_identity(
        self, session: AsyncSession, user_auth_repo: UserAuth
    ):
        user = await _create_user(session, username="bob")
        await _create_identity(
            user_auth_repo, user_id=user.id, provider="github", subject="bob_github"
        )

        maybe_user = await user_auth_repo.get_user_by_identity(
            provider="github", subject="bob_github"
        )
        got = maybe_user.value_or(None)
        assert got is not None
        assert got.id == user.id

    async def test_get_auth_identity_by_user_id_and_provider(
        self, session: AsyncSession, user_auth_repo: UserAuth
    ):
        user = await _create_user(session)
        await _create_identity(
            user_auth_repo, user_id=user.id, provider="google", subject="uid-123"
        )

        maybe_identity = await user_auth_repo.get_auth_identity_by_user_id_and_provider(
            user_id=user.id, provider="google"
        )
        identity = maybe_identity.value_or(None)
        assert identity is not None
        assert identity.user_id == user.id
        assert identity.provider == "google"
        assert identity.subject == "uid-123"

    async def test_touch_last_login_by_identity(
        self, session: AsyncSession, user_auth_repo: UserAuth
    ):
        user = await _create_user(session)
        created = await _create_identity(
            user_auth_repo, user_id=user.id, provider="apple", subject="sub-apple"
        )
        # Depending on your model default, adjust this assertion if needed:
        assert created.last_login_at is None

        res = await user_auth_repo.touch_last_login_by_identity(
            provider="apple", subject="sub-apple"
        )
        assert isinstance(res, Success)
        updated = res.unwrap()
        assert updated.last_login_at is not None
        assert isinstance(updated.last_login_at, datetime)

    async def test_update_auth_identity_empty_values(
        self, session: AsyncSession, user_auth_repo: UserAuth
    ):
        user = await _create_user(session)
        await _create_identity(
            user_auth_repo, user_id=user.id, provider="password", subject="to-update"
        )

        res = await user_auth_repo.update_auth_identity(
            user_id=user.id, provider="password", values={}
        )
        # Expect Failure(ValueError("No fields provided to update."))
        assert isinstance(res, Failure)

    async def test_update_auth_identity_success(
        self, session: AsyncSession, user_auth_repo: UserAuth
    ):
        user = await _create_user(session)
        await _create_identity(
            user_auth_repo,
            user_id=user.id,
            provider="password",
            subject="to-update-ok",
            email="old@example.com",
            email_verified=False,
        )

        res = await user_auth_repo.update_auth_identity(
            user_id=user.id,
            provider="password",
            values={
                "email": "new@example.com",
                "email_verified": True,
                "hash_method": "argon2",
            },
        )
        assert isinstance(res, Success)
        updated = res.unwrap()
        assert updated.email == "new@example.com"
        assert updated.email_verified is True
        assert updated.hash_method == "argon2"

        # Verify persisted changes
        row = (
            await session.execute(
                select(AuthIdentities).where(AuthIdentities.id == updated.id)
            )
        ).scalar_one_or_none()
        assert row is not None
        assert row.email == "new@example.com"
        assert row.email_verified is True
        assert row.hash_method == "argon2"

    async def test_delete_auth_identity(
        self, session: AsyncSession, user_auth_repo: UserAuth
    ):
        user = await _create_user(session)
        identity = await _create_identity(
            user_auth_repo, user_id=user.id, provider="github", subject="to-delete"
        )

        res = await user_auth_repo.delete_auth_identity(identity.id)
        assert isinstance(res, Success)
        assert res.unwrap() == 1  # Number of rows deleted

        # Verify it is gone
        row = (
            await session.execute(
                select(AuthIdentities).where(AuthIdentities.id == identity.id)
            )
        ).scalar_one_or_none()
        assert row is None
