import os


class Config:
    DEBUG_MODE = os.getenv("DEBUG_MODE", "False").lower() == "true"
    DATABASE_URL = os.getenv(
        "DATABASE_URL",
        "postgresql+asyncpg://postgres:secret@localhost:5432/postgres",
    )
