# HWatchDog Scraper

Scraps amd / intel cpu prices from online ecommerce websites.

### Support

Supported sites
- Amazon { USA, CA}

Upcoming
- Newegg
- Microcenter
- Canada Computers 

### DB Requirements

- Rename ``copy.env`` -> ``.env``
- Fill ``.env`` with db information.
- Run DDL contained in ``DDL.sql`` on your database. 

### DOCKER

Build Image: ``docker build -t hwatchdog-s .``

Run Container: ``docker run -d -p 9001:9001 --restart=always hwatchdog-s`` on port 9001

You can access the Hwatchdog ui on port 9001




