import pytest

# TODO: more tests here


def test_post_creation(make_user, auth_user, make_request, core_addr):
    user = auth_user(make_user())
    resp = make_request(
        'POST', core_addr, '/v1/posts', data={'content': 'post_content'}, cookies=user.auth_cookies
    )
    assert resp.status_code == 200
    assert 'post_id' in resp.json()


def test_post_retrieve(make_user, auth_user, make_request, core_addr):
    user = auth_user(make_user())
    post_content = 'post_content'
    resp = make_request(
        'POST', core_addr, '/v1/posts', data={'content': post_content}, cookies=user.auth_cookies
    )
    assert resp.status_code == 200
    post_id = resp.json()['post_id']

    resp = make_request('GET', core_addr, f'/v1/posts/{post_id}', cookies=user.auth_cookies)
    assert resp.status_code == 200
    post = resp.json()
    assert post['id'] == post_id
    assert post['content'] == post_content
    assert 'author_id' in post   # TODO: test author_id is valid
