import time
import uuid
from datetime import datetime

import pytest
from returns.result import Failure, Success
from sqlalchemy.ext.asyncio import AsyncSession

from annotatex.models.user import User, UserPortfolios
from annotatex.repositories.user_portfolio_repository import UserPortfolioRepository


@pytest.fixture()
def portfolio_repo(session: AsyncSession) -> UserPortfolioRepository:
    return UserPortfolioRepository(session)


async def _create_user(session: AsyncSession, *, username: str = "u1") -> User:
    user = User(username=username, is_active=True)
    async with session.begin():
        session.add(user)
        await session.flush()
    return user


async def _create_portfolio(
    repo: UserPortfolioRepository, *, user_id: uuid.UUID, display_name: str
) -> UserPortfolios:
    result = await repo.create_portfolio(user_id=user_id, display_name=display_name)
    assert isinstance(result, Success)
    return result.unwrap()


class TestUserPortfolioRepository:
    async def test_create_portfolio_success(
        self, portfolio_repo: UserPortfolioRepository, session: AsyncSession
    ):
        user = await _create_user(session, username="alice")
        result = await portfolio_repo.create_portfolio(
            user_id=user.id, display_name="Core"
        )
        assert isinstance(result, Success)

        pf = result.unwrap()
        assert isinstance(pf.id, int)
        assert pf.user_id == user.id
        assert pf.displayName == "Core"

    async def test_update_portfolio_success(
        self, portfolio_repo: UserPortfolioRepository, session: AsyncSession
    ):
        user = await _create_user(session, username="charlie")
        pf = await _create_portfolio(
            portfolio_repo, user_id=user.id, display_name="Old"
        )

        before_updated_at = getattr(pf, "updated_at", None)
        # Sleep for a second for different time.
        time.sleep(1)
        upd = await portfolio_repo.update_portfolio(
            pf.id,
            {"displayName": "New"},
        )
        assert isinstance(upd, Success)
        new_pf = upd.unwrap()
        assert new_pf.displayName == "New"

        if before_updated_at is not None:
            assert isinstance(new_pf.updated_at, datetime)
            assert new_pf.updated_at != before_updated_at

    async def test_update_portfolio_not_found(
        self, portfolio_repo: UserPortfolioRepository
    ):
        bad = await portfolio_repo.update_portfolio(
            portfolio_id=999999,
            values={"displayName": "X"},
        )
        assert isinstance(bad, Failure)

    async def test_update_portfolio_no_fields(
        self, portfolio_repo: UserPortfolioRepository, session: AsyncSession
    ):
        user = await _create_user(session, username="dave")
        pf = await _create_portfolio(
            portfolio_repo, user_id=user.id, display_name="Name"
        )

        no_fields = await portfolio_repo.update_portfolio(pf.id, {})
        assert isinstance(no_fields, Failure)
