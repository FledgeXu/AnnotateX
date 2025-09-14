from pydantic import BaseModel

from annotatex.models.organization import OrganizationKind


class CreateOrganizationSchema(BaseModel):
    name: str
    kind: OrganizationKind
