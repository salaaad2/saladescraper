import argparse
from urllib.request import urlopen

def init_argparse() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        usage="%(prog)s [OPTION] [TARGET_NAME]...",
        description="get all articles from a substack"
    )

    parser.add_argument(
        "-v", "--version", action="version",
        version = f"{parser.prog} version 0.0.1"
    )
    parser.add_argument('target', nargs='*')
    return parser

def getArticles(target_name) -> str:
    full_target = "{}{}{}{}{}".format(
        "https://",
        target_name,
        ".substack.com",
    )
    print("saladescraper: targeting: {}", full_target)
    post_urls = []
    home_request = urllib2.urlopen(full_target).read()
    home_content = home_request.read()

    # repeating pattern
    post_str = full_target + "/p/"
    post_start_index = str.find(post_str)
    post_end_index = str.find("/comments")
    post_count = 0
    while post_start_index != -1 and post_end_index != -1:
        if (post_end_index - post_start_index) > 100:
            post_str = post_str[post_start_index + len(post_str)]
            post_start_index = str.find(post_str)
            post_end_index = str.find("/comments")
        else:
            post_str = post_str[post_start_index + len(post_str + "/comments")]
            post_start_index = str.find(post_str)
            post_end_index = str.find("/comments")
            post_count += 1
            post_urls.push_back(post_str.split(post_start_index, post_end_index - post_start_index))
    if post_count == 0:
        return "no posts found"

    home_request.close()
    if len(main_contents) == 0:
        return "error"

    for post_url in post_urls:
        print(post_url)
    return "success"

def main() -> None:
    arg_parser = init_argparse()
    args = arg_parser.parse_args()
    getArticles(args.target)
