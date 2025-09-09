import uuid
from datetime import datetime
from enum import Enum

from sqlalchemy import Boolean, DateTime, String, Uuid, func
from sqlalchemy import Enum as SqlEnum
from sqlalchemy.orm import Mapped, mapped_column

from annotatex.db.base import Base


class OrganizationKind(Enum):
    client = "client"
    vendor = "vendor"
    internal = "internal"


class Organization(Base):
    __tablename__ = "organization"

    id: Mapped[uuid.UUID] = mapped_column(Uuid, primary_key=True, default=uuid.uuid4)
    name: Mapped[str] = mapped_column(String, unique=True, nullable=False)
    kind: Mapped[OrganizationKind] = mapped_column(
        SqlEnum(OrganizationKind, name="organization_kind"), nullable=False
    )

    is_active: Mapped[bool] = mapped_column(Boolean, default=True)
    created_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), nullable=False, server_default=func.now()
    )
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        nullable=False,
        server_default=func.now(),
        onupdate=func.now(),
    )
