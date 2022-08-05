import argparse
import json
import requests
import threading
import os

def init_argparse() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        usage="%(prog)s [OPTION] [TARGET_NAME]...",
        description="get all articles from a substack"
    )

    parser.add_argument(
        "-v", "--version", action="version",
        version = f"{parser.prog} version 0.0.1"
    )
    parser.add_argument('target', type=str)
    return parser

def get_article(article_url):
    article_body = requests.get(article_url).text
    art_name = article_url.split('/')[-1]
    article_path = "./{}".format(
        art_name
    )
    fo = open(art_name, "w")
    fo.write(article_body)
    fo.close()

def get_article_list(target_name):
    full_target = "{}{}{}{}".format(
        "https://",
        target_name,
        ".substack.com",
        "/api/v1/archive?sort=new&search=&offset=1&limit=100",
    )
    print(f"saladescraper: targeting: {full_target}")
    post_count = 0
    post_urls = []
    home_content = requests.get(full_target).text

    json_response = json.loads(home_content)
    post_list = []
    for item in json_response:
        post_urls.append(item['canonical_url'])
        post_count += 1

    if post_count == 0:
        return ["no posts found", ""]
    else:
        print(post_count)
    return post_urls

def main():
    arg_parser = init_argparse()
    args = arg_parser.parse_args()
    articles_list = get_article_list(args.target)
    thread_list = []

    if not os.path.exists(args.target):
        os.mkdir(args.target)
    os.chdir(args.target)
    thread_list = []
    for article in articles_list:
        t = threading.Thread(target=get_article, args=(article,))
        thread_list.append(t)

    for t in thread_list:
        t.start()
    for t in thread_list:
        t.join()
    return 0

if __name__ == '__main__':
    main()
