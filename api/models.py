from datetime import datetime

from sqlmodel import Field, Relationship, Session, SQLModel, create_engine, select


class UserBase(SQLModel):
    username: str = Field(index=True)


class User(UserBase, table=True):
    id: int | None = Field(default=None, primary_key=True)
    hashed_password: str = Field()
    created_at: datetime

class UserCreate(UserBase):
    password: str
    created_at: datetime = datetime.now()

class UserPublic(UserBase):
    id: int
