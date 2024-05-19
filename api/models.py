
from datetime import date, datetime

from sqlmodel import Field, Relationship, SQLModel

''' User Models '''
class UserBase(SQLModel):
    username: str = Field(index=True)
    email: str = Field(index=True)


class User(UserBase, table=True):
    id: int | None = Field(default=None, primary_key=True)
    hashed_password: str = Field()
    created_at: datetime
    updated_at: datetime

    items: list["Item"] = Relationship(back_populates="user")
    assets: list["Asset"] = Relationship(back_populates="user")

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
    email: str | None = None
    updated_at: datetime = datetime.now()

''' Plaid API Models '''
# Link Events Table logs responses from Plaid API for client requests to the Plaid Link client
class LinkEvent(SQLModel, table=True):
    id: int | None = Field(default=None, primary_key=True)
    type: str
    user_id: int
    link_session_id: int
    request_id: str = Field(unique=True)
    error_type: str
    error_code: str
    status: str
    created_at: datetime = datetime.now()

# API Events table logs responses from Plaid API for server requests to the Plaid client
class APIEvent(SQLModel, table=True):
    id: int | None = Field(default=None, primary_key=True)
    item_id: int
    user_id: int
    plaid_method: str
    arguments: str | None
    request_id: str = Field(unique=True)
    error_type: str
    error_code: str
    created_at: datetime = datetime.now()

# Each Item represents a log-in for a financial institution for Plaid (if user has 2 accounts with 1 institution, they wil be under 1 Item)
class Item(SQLModel, table=True):
    id: int | None = Field(default=None, primary_key=True)
    plaid_access_token: str = Field(unique=True)
    plaid_item_id: str = Field(unique=True)
    plaid_institution_id: str
    status: str
    created_at: datetime = datetime.now()
    updated_at: datetime = datetime.now()
    transactions_cursor: str
    # cursor keeps track of most recent transactions update

    user_id: int = Field(foreign_key="user.id")
    user: User = Relationship(back_populates='items')

    accounts: list["Account"] = Relationship(back_populates="item")


class Asset(SQLModel, table=True):
    id: int | None = Field(default=None, primary_key=True)
    value: float
    description: str | None
    created_at: datetime = datetime.now()
    updated_at: datetime = datetime.now()

    user_id: int = Field(foreign_key="user.id")
    user: User = Relationship(back_populates="assets")


class Account(SQLModel, table=True):
    id: int | None = Field(default=None, primary_key=True)
    plaid_account_id: str = Field(unique=True, index=True)
    name: str = Field(index=True)
    mask: str
    official_name: str | None
    available_balance: float
    iso_currency_code: str | None
    unofficial_currency_code: str | None
    type: str
    subtype: str
    created_at: datetime = datetime.now()
    updated_at: datetime = datetime.now()

    item_id: int = Field(foreign_key="item.id")
    item: Item = Relationship(back_populates="accounts")

    transactions: list["Transaction"] = Relationship(back_populates="account")


class Transaction(SQLModel, table=True):
    id: int | None = Field(default=None, primary_key=True)
    plaid_transaction_id: str = Field(unique=True, index=True)
    plaid_category_id: str | None
    category: str | None
    subcategory: str | None
    transaction_type: str
    name: str
    amount: float
    iso_currency_code: str | None
    unofficial_currency_code: str | None
    date: date
    pending: bool
    account_owner: str | None
    created_at: datetime = datetime.now()
    updated_at: datetime = datetime.now()

    account_id: int = Field(foreign_key="account.id")
    account: Account = Relationship(back_populates="transactions")
