
import pytest
import requests
import uuid
import time

BASE_URL = "http://127.0.0.1:8000"


def create_task():
    """Создаёт задачу и возвращает её ID."""
    response = requests.post(
        f"{BASE_URL}/task",
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


def test_create_task():
    task_id = create_task()
    assert uuid.UUID(task_id)


def test_task_status_and_result():
    task_id = create_task()

    status_url = f"{BASE_URL}/status/{task_id}"
    result_url = f"{BASE_URL}/result/{task_id}"

    deadline = time.monotonic() + 10

    while time.monotonic() < deadline:
        response = requests.get(status_url, timeout=5)

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

    response = requests.get(result_url, timeout=5)

    assert response.status_code == 200

    data = response.json()
    assert "result" in data
    assert data["result"] != ""


def test_task_not_found():
    invalid_task_id = str(uuid.uuid4())

    status_url = f"{BASE_URL}/status/{invalid_task_id}"
    result_url = f"{BASE_URL}/result/{invalid_task_id}"

    response = requests.get(status_url, timeout=5)
    assert response.status_code == 404

    response = requests.get(result_url, timeout=5)
    assert response.status_code == 404


def test_task_requires_payload():
    response = requests.post(
        f"{BASE_URL}/task",
        timeout=5
    )

    assert response.status_code == 400
