from contextlib import asynccontextmanager

from sqlmodel import Session, SQLModel, select
from fastapi import Depends, FastAPI, HTTPException, Query
from passlib.hash import pbkdf2_sha256

import models
from engine import engine

def create_db_and_tables():
    SQLModel.metadata.create_all(engine)

@asynccontextmanager
async def lifespan(app: FastAPI):
    create_db_and_tables()
    yield

app = FastAPI(lifespan=lifespan)

def get_session():
    with Session(engine) as session:
        yield session


@app.post("/users/", response_model=models.UserPublic)
def create_user(*, session: Session = Depends(get_session), user: models.UserCreate):
    hashed_password = pbkdf2_sha256.hash(user.password)
    extra_data = {"hashed_password": hashed_password}
    db_user = models.User.model_validate(user, update=extra_data)
    session.add(db_user)
    session.commit()
    session.refresh(db_user)
    return db_user

@app.get("/users/", response_model=list[models.UserPublic])
def read_users(session: Session = Depends(get_session), offset: int = 0, limit: int = Query(default=100, le=100)):
    return session.exec(select(models.User).offset(offset).limit(limit)).all()
