
from datetime import datetime

from sqlmodel import Field, Relationship, SQLModel

# User Models
class UserBase(SQLModel):
    username: str = Field(index=True)
    email: str = Field(index=True)


class User(UserBase, table=True):
    id: int | None = Field(default=None, primary_key=True)
    hashed_password: str = Field()
    created_at: datetime
    updated_at: datetime

class UserCreate(UserBase):
    password: str
    created_at: datetime = datetime.now()
    updated_at: datetime = datetime.now()

class UserPublic(UserBase):
    id: int
    username: str
    created_at: datetime
    updated_at: datetime

class UserUpdate(UserBase):
    username: str | None = None
    password: str | None = None
    updated_at: datetime = datetime.now()
