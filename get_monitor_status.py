import os
import sys
import requests
import time

def get_monitor_id_by_name(api_key, app_key, monitor_name):
    base_url = "https://api.datadoghq.com/api/v1/monitor"
    headers = {
        "Content-Type": "application/json",
        "DD-API-KEY": api_key,
        "DD-APPLICATION-KEY": app_key,
    }
    params = {
        "query": f"name:{monitor_name}",
    }

    try:
        response = requests.get(base_url, headers=headers, params=params)

        if response.status_code == 200:
            monitors_data = response.json()
            for monitor in monitors_data:
                if monitor["name"] == monitor_name:
                    return monitor["id"]

            # If the monitor with the given name was not found
            print(f"Monitor with name '{monitor_name}' not found.")
            return None

        else:
            print(f"Failed to fetch monitors. Status code: {response.status_code}")
            return None

    except requests.exceptions.RequestException as e:
        print(f"Error: {e}")
        return None

def get_monitor_status(api_key, app_key, monitor_id):
    base_url = f"https://api.datadoghq.com/api/v1/monitor/{monitor_id}"
    headers = {
        "Content-Type": "application/json",
        "DD-API-KEY": api_key,
        "DD-APPLICATION-KEY": app_key,
    }

    try:
        response = requests.get(base_url, headers=headers)

        if response.status_code == 200:
            monitor_data = response.json()
            return monitor_data["overall_state"]
        else:
            print(f"Failed to get monitor status. Status code: {response.status_code}")
            return None

    except requests.exceptions.RequestException as e:
        print(f"Error: {e}")
        return None

if __name__ == "__main__":
    monitor_name = sys.argv[1]
    api_key = os.environ["API_KEY"]
    app_key = os.environ["APP_KEY"]

    # Get the monitor ID by name
    monitor_id = get_monitor_id_by_name(api_key, app_key, monitor_name)

    if monitor_id is not None:
        status = get_monitor_status(api_key, app_key, monitor_id)
        print(f"Monitor ID {monitor_name} - Status: {status}")
