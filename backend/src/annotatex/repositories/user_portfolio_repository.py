import uuid
from typing import Any, Mapping

from returns.result import Failure, Result, Success
from sqlalchemy import func, update

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

    async def update_portfolio(
        self,
        portfolio_id: int,
        values: Mapping[str, Any],
    ) -> Result[UserPortfolios, Exception]:
        if not values:
            return Failure(ValueError("No fields provided to update."))

        payload = dict(values)
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
