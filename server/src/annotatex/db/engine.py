from typing import AsyncGenerator

from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

from annotatex.config import Config

engine = create_async_engine(
    Config.DATABASE_URL,
    echo=Config.DEBUG_MODE,
    pool_pre_ping=True,
)

AsyncSessionLocal = async_sessionmaker(
    engine, expire_on_commit=False, autoflush=False, class_=AsyncSession
)


async def get_async_session() -> AsyncGenerator[AsyncSession, None]:
    async with AsyncSessionLocal() as session:
        yield session
