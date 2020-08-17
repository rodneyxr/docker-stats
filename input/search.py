#!/usr/bin/env python3
import argparse
import base64
import os
import time

from github import Github
from github.ContentFile import ContentFile


def main():
    # Handle Arguments
    parser = argparse.ArgumentParser()
    parser.add_argument('--token', required=True, help="your GitHub token")
    parser.add_argument('--count', default=1000, type=int, help="number of files to download")
    args = parser.parse_args()

    # Authenticate with GitHub API and authenticate
    g = Github(args.token)
    results = g.search_code(query="filename:install.sh")

    # Iterate through each file until download count has been reached
    count = args.count
    total_downloaded = 0
    for file in results:
        file: ContentFile
        if total_downloaded < count:
            download_file(file)
            total_downloaded += 1
        else:
            return
        time.sleep(.1)

        # Wait for rate limit to reset before continuing, if needed
        rate_limit = g.get_rate_limit()
        if rate_limit.search.remaining == 0:
            while rate_limit.search.remaining < 0:
                print("waiting for rate_limit.search to reset...")
                time.sleep(5)

    # Ensure there are no duplicates in our sources.txt
    clean_sources()


def clean_sources():
    with open("sources.txt") as f:
        sources = set(f.readlines())
    with open("sources.txt", "w") as f:
        f.writelines(sources)


def download_file(content_file):
    """
    Saves the content_file to a file on disk.
    :param content_file: The ContentFile object to save
    :return: True if the file was saved.
    """
    filename = f'{content_file.repository.owner.login}_{content_file.repository.name}_{content_file.name}'
    filepath = os.path.join("searches", filename)
    os.makedirs(os.path.dirname(filepath), exist_ok=True)
    print(filepath)
    with open("sources.txt", "a+") as f:
        f.write(content_file.download_url + "\n")
    if not os.path.exists(filepath):
        with open(filepath, "wb") as f:
            content = base64.b64decode(content_file.content)
            f.write(content)
        return True


if __name__ == "__main__":
    main()
