
import subprocess
import csv
import sys
import argparse
import os
import json

durl = "https://api.github.inl.gov/repos/CyOTE/CATCH/import/issues"

class IssuesPrompt:
    def __init__(self, url=None, csv_r=None, token=None):
        self.url = url
        self.csv_r=csv_r
        self.token=token
        self.title=""
        self.body=""
        self.labels=[]
        self.assignee=""
        self.milestone=""
        self.closed=False
        self.extra=""

    def post(self):
        # I haven't figured out milestones yet
        # data_inner =f"\"title\":\"{self.title}\",\"body\":\"{self.body}\",\"labels\":{self.labels},\"assignee\":\"{self.assignee}\", \"milestone\":\"{self.milestone}\""
        data_inner =f"\"title\":\"{self.title}\",\"body\":\"{self.body}\",\"labels\":{self.labels},\"assignee\":\"{self.assignee}\"" 

        if self.closed:
            data_inner = data_inner + f",\"closed\":true"
        data_outer = '{"issue":{' + data_inner + '} }'
        print(data_outer)
        with subprocess.Popen(["curl", "--request", "POST", "--url",
        self.url, "--header", f"Authorization: Bearer {self.token}", 
        "--header","X-GitHub-Api-Version: 2022-11-28", "-d", f"{data_outer}"
        ], stderr=subprocess.STDOUT) as a:
            print(a.stdout)

    def prompt(self):
        print("Enter Fields for Issue:\n")
        self.title=input("Title: ")
        self.body=input("Body: ")
        labels = input("Labels(comma sep, no spaces): ")
        self.labels = json.dumps(labels.split(","))
        print(self.labels) # -wb because I'm afraid formatting will fail
        self.assignee=input("Assignee: ")
        self.extra=input("Extra Options(must be formatted): ")

        self.post()

    def parse(self):
        if not os.path.exists(self.csv_r):
            print("csv does not exist")
            return
        with open(self.csv_r, 'r') as csv_f:
            reader = csv.reader(csv_f)
            for row in reader:
                if row == ['Title', 'Description', 'Assignee', 'State', 'Milestone', 'Labels']:
                    continue #skip first line
             
                self.title=row[0]
                self.body=row[1].replace('"','').replace('\'', '').replace("\n", "\\n")
                self.assignee=row[2]
                self.milestone=row[4]
                if row[5] != "":
                    self.labels=json.dumps(row[5].split(","))
                else:
                    self.labels=[]
                if "Closed" in row[3]:
                    self.closed=True
                else:
                    self.closed=False

                if self.body == "":
                    self.body = "-blank-"
                self.post()

    def setup(self):
        if self.url or self.token is None:
            print(f"No url given, default: {durl}")
        if self.csv_r is None:
            self.prompt()
            return
        self.parse()

    
        
def parse_args(args):
    parser = argparse.ArgumentParser(
        description="Posts issues to github"
    )
    parser.add_argument("--url", nargs='?', type=str, action="store", default=durl)
    parser.add_argument("--token",nargs='?', type=str, action="store", default=None)
    parser.add_argument("--csv", nargs='?',type=str, action="store", default=None)

    return parser.parse_args(args)
    

def main():
    args = parse_args(sys.argv[1:])
    iss = IssuesPrompt(args.url, args.csv, args.token)
    iss.setup()

if __name__ == "__main__":
    main()

