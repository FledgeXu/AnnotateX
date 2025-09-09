from uuid import UUID

from returns.maybe import Maybe
from sqlalchemy import select

from annotatex.models.user import User
from annotatex.repositories.base_repository import BaseRepository


class UserRepository(BaseRepository):
    async def get_user_by_id(self, user_id: UUID) -> Maybe[User]:
        async with self.session.begin():
            result = await self.session.execute(select(User).where(User.id == user_id))
            return Maybe.from_optional(result.scalar_one_or_none())

    async def get_user_by_username(self, username: str) -> Maybe[User]:
        async with self.session.begin():
            result = await self.session.execute(
                select(User).where(User.username == username)
            )
            return Maybe.from_optional(result.scalar_one_or_none())
