import uuid
from typing import Any, Mapping, Sequence

from returns.result import Failure, Result, Success
from sqlalchemy import delete, func, select, update

from annotatex.models.user import UserPortfolios
from annotatex.repositories.base_repository import BaseRepository


class UserPortfolioRepository(BaseRepository):
    async def create_portfolio(
        self, *, user_id: uuid.UUID, display_name: str
    ) -> Result[UserPortfolios, Exception]:
        portfolio = UserPortfolios(user_id=user_id, displayName=display_name)
        try:
            async with self.session.begin():
                self.session.add(portfolio)
                await self.session.flush()
            return Success(portfolio)
        except Exception as exc:
            return Failure(exc)

    async def list_portfolios(self, user_id: uuid.UUID) -> Sequence[UserPortfolios]:
        result = await self.session.execute(
            select(UserPortfolios)
            .where(UserPortfolios.user_id == user_id)
            .order_by(UserPortfolios.id.asc())
        )
        return list(result.scalars().all())

    async def delete_portfolio(self, portfolio_id: int) -> Result[int, Exception]:
        try:
            async with self.session.begin():
                result = await self.session.execute(
                    delete(UserPortfolios).where(UserPortfolios.id == portfolio_id)
                )
                rows = result.rowcount or 0
            return Success(rows)
        except Exception as exc:
            return Failure(exc)

    async def update_portfolio(
        self,
        portfolio_id: int,
        values: Mapping[str, Any],
        *,
        touch_updated_at: bool = False,
    ) -> Result[UserPortfolios, Exception]:
        if not values and not touch_updated_at:
            return Failure(ValueError("No fields provided to update."))

        payload = dict(values)
        if touch_updated_at:
            payload["updated_at"] = func.now()

        try:
            async with self.session.begin():
                stmt = (
                    update(UserPortfolios)
                    .where(UserPortfolios.id == portfolio_id)
                    .values(**payload)
                    .returning(UserPortfolios)
                )
                result = await self.session.execute(stmt)
                updated = result.scalar_one_or_none()
                if updated is None:
                    raise ValueError("Not found")
            return Success(updated)
        except Exception as exc:
            return Failure(exc)
