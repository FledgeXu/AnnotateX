import os

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine
from sqlalchemy.orm import close_all_sessions

from annotatex.db.base import Base


def _get_test_db_url() -> str:
    return os.getenv("TEST_DATABASE_URL", "sqlite+aiosqlite:///:memory:")


@pytest.fixture(scope="function")
async def engine():
    url = _get_test_db_url()
    engine = create_async_engine(url, future=True)
    try:
        yield engine
    finally:
        await engine.dispose()
        close_all_sessions()


@pytest.fixture()
async def session(engine):
    maker = async_sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.drop_all)
        await conn.run_sync(Base.metadata.create_all)

        s: AsyncSession = maker(bind=conn)
        try:
            yield s
        finally:
            await s.close()
