import time
import uuid

import pytest
import requests

BASE_URL = "http://127.0.0.1:8000"


@pytest.fixture(scope="module")
def auth_headers():
    credentials = {
        "username": f"user_{uuid.uuid4()}",
        "password": "test-password",
    }

    response = requests.post(f"{BASE_URL}/register", json=credentials, timeout=5)
    assert response.status_code == 201, response.text

    response = requests.post(f"{BASE_URL}/login", json=credentials, timeout=5)
    assert response.status_code == 200, response.text
    token = response.json()["token"]
    assert isinstance(token, str) and token

    return {"Authorization": f"Bearer {token}"}


def create_task(auth_headers):
    """Создаёт задачу и возвращает её ID."""
    response = requests.post(
        f"{BASE_URL}/task",
        headers=auth_headers,
        json={
            "code": "print(1 + 1)",
            "language": "python"
        },
        timeout=5
    )

    assert response.status_code == 201

    data = response.json()
    assert "task_id" in data

    return data["task_id"]


def test_create_task(auth_headers):
    task_id = create_task(auth_headers)
    assert uuid.UUID(task_id)


def test_task_status_and_result(auth_headers):
    task_id = create_task(auth_headers)

    status_url = f"{BASE_URL}/status/{task_id}"
    result_url = f"{BASE_URL}/result/{task_id}"

    deadline = time.monotonic() + 10

    while time.monotonic() < deadline:
        response = requests.get(status_url, headers=auth_headers, timeout=5)

        assert response.status_code == 200

        data = response.json()
        assert "status" in data

        if data["status"] == "ready":
            break

        assert data["status"] == "in_progress", (
            f"undefined status: {data['status']}!"
        )

        time.sleep(1)
    else:
        pytest.fail("Task processing timeout")

    response = requests.get(result_url, headers=auth_headers, timeout=5)

    assert response.status_code == 200

    data = response.json()
    assert "result" in data
    assert data["result"] != ""


def test_task_not_found(auth_headers):
    invalid_task_id = str(uuid.uuid4())

    status_url = f"{BASE_URL}/status/{invalid_task_id}"
    result_url = f"{BASE_URL}/result/{invalid_task_id}"

    response = requests.get(status_url, headers=auth_headers, timeout=5)
    assert response.status_code == 404

    response = requests.get(result_url, headers=auth_headers, timeout=5)
    assert response.status_code == 404


def test_task_requires_payload(auth_headers):
    response = requests.post(
        f"{BASE_URL}/task",
        headers=auth_headers,
        timeout=5
    )

    assert response.status_code == 400


@pytest.mark.parametrize(
    "headers",
    [{}, {"Authorization": "Bearer not-a-session-token"}],
    ids=["missing-token", "unknown-token"],
)
def test_task_endpoints_require_session(headers):
    task_id = uuid.uuid4()
    endpoints = (
        ("POST", "/task"),
        ("GET", f"/status/{task_id}"),
        ("GET", f"/result/{task_id}"),
    )

    for method, path in endpoints:
        response = requests.request(method, f"{BASE_URL}{path}", headers=headers, timeout=5)
        assert response.status_code == 401, f"{method} {path}: {response.status_code}"
