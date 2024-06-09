import json
import typing
import dataclasses
import uuid
import string
import random

import pytest
import requests


@pytest.fixture(name='make_request')
def _make_request():
    def _impl(method, addr, handle, params=None, data=None, cookies=None):
        if data is not None:
            data = json.dumps(data)
        req = requests.Request(
            method, 'http://' + addr + handle, params=params, data=data, cookies=cookies
        )
        prepared = req.prepare()
        s = requests.Session()
        resp = s.send(prepared)
        return resp

    return _impl


@dataclasses.dataclass
class User:
    login: str
    password: str
    auth_cookies: typing.Optional[dict] = None


@pytest.fixture(name='core_addr')
def _core_addr():
    return 'core_service:3000'


@pytest.fixture(name='make_user')
def _make_user(make_request, core_addr):
    def _impl() -> User:
        login = ''.join(random.sample(string.ascii_letters, 12))
        password = 'val1dpassw0rd'
        resp = make_request(
            'POST', core_addr, '/v1/users', data={'login': login, 'password': password}
        )
        assert resp.status_code == 200
        return User(login, password)

    return _impl


@pytest.fixture(name='auth_user')
def _auth_user(make_request, core_addr):
    def _impl(user: User) -> User:
        resp = make_request(
            'POST', core_addr, '/v1/auth', data={'login': user.login, 'password': user.password}
        )
        assert resp.status_code == 200
        user.auth_cookies = resp.cookies.get_dict()
        return user

    return _impl
