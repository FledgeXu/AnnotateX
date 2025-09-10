import uuid
from typing import Any, Mapping

from returns.maybe import Maybe
from returns.result import Failure, Result, Success
from sqlalchemy import delete, func, select, update
from sqlalchemy.exc import IntegrityError

from annotatex.models.user import AuthIdentities, User
from annotatex.repositories.base_repository import BaseRepository


class UserAuth(BaseRepository):
    async def add_auth_identity(
        self,
        *,
        user_id: uuid.UUID,
        provider: str,
        subject: str,
        email: str | None = None,
        email_verified: bool = False,
        password_hash: str | None = None,
    ) -> Result[AuthIdentities, IntegrityError]:
        identity = AuthIdentities(
            user_id=user_id,
            provider=provider,
            subject=subject,
            email=email,
            email_verified=email_verified,
            password_hash=password_hash,
        )
        try:
            async with self.session.begin():
                self.session.add(identity)
                await self.session.flush()
            return Success(identity)
        except IntegrityError as exc:
            return Failure(exc)

    async def get_user_by_identity(self, *, provider: str, subject: str) -> Maybe[User]:
        result = await self.session.execute(
            select(User)
            .join(AuthIdentities, AuthIdentities.user_id == User.id)
            .where(
                AuthIdentities.provider == provider,
                AuthIdentities.subject == subject,
            )
        )
        return Maybe.from_optional(result.scalar_one_or_none())

    async def touch_last_login_by_identity(
        self, *, provider: str, subject: str
    ) -> Result[AuthIdentities, Exception]:
        try:
            async with self.session.begin():
                stmt = (
                    update(AuthIdentities)
                    .where(
                        AuthIdentities.provider == provider,
                        AuthIdentities.subject == subject,
                    )
                    .values(last_login_at=func.now())
                    .returning(AuthIdentities)
                )
                result = await self.session.execute(stmt)
                row = result.scalar_one_or_none()
                if row is None:
                    raise ValueError("Auth identity not found")
            return Success(row)
        except Exception as exc:
            return Failure(exc)

    async def update_auth_identity(
        self, *, user_id: uuid.UUID, provider: str, values: Mapping[str, Any]
    ) -> Result[AuthIdentities, Exception]:
        if not values:
            return Failure(ValueError("No fields provided to update."))

        try:
            async with self.session.begin():
                result = await self.session.execute(
                    select(AuthIdentities).where(
                        AuthIdentities.user_id == user_id,
                        AuthIdentities.provider == provider,
                    )
                )
                existing = result.scalar_one_or_none()

                if existing is None:
                    raise ValueError("Not found")

                values = {"updated_at": func.now(), **values}
                stmt = (
                    update(AuthIdentities)
                    .where(AuthIdentities.id == existing.id)
                    .values(**values)
                    .returning(AuthIdentities)
                )
                result = await self.session.execute(stmt)
                updated = result.scalar_one()

            return Success(updated)
        except Exception as exc:
            return Failure(exc)

    async def delete_auth_identity(self, identity_id: int) -> Result[int, Exception]:
        try:
            async with self.session.begin():
                result = await self.session.execute(
                    delete(AuthIdentities).where(AuthIdentities.id == identity_id)
                )
                rows = result.rowcount or 0
            return Success(rows)
        except Exception as exc:
            return Failure(exc)
