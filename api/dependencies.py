from contextlib import asynccontextmanager

from sqlmodel import Session, SQLModel
from fastapi import FastAPI

from .engine import engine

def create_db_and_tables():
    SQLModel.metadata.create_all(engine)

def get_session():
    with Session(engine) as session:
        yield session

@asynccontextmanager
async def lifespan(app: FastAPI):
    create_db_and_tables()
    yield
