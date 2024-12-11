# Project Managment Demo Backend

[![dhtmlx.com](https://img.shields.io/badge/made%20by-DHTMLX-blue)](https://dhtmlx.com/)

[How to start](#how-to-start) | [API](#api) | [License](#license) | [Useful links](#links) | [Join our online community](#join)

<a name="how-to-start"></a>
## How to start

Using docker-compose:

```
docker-compose up --build
```

<a name="api"></a>
## API

### REST API

Each widget has own REST API in the following format:

#### backend.url.com/api/{widget}/routes

See more details for each widget api in:

- api/kanban.go
- api/gantt.go
- api/todo.go
- api/scheduler.go

### WS API

Each widget has own Web Socket API in the following format:

#### backend.url.com/api/{widget}/v1

See more details for each widget in:

- publisher/kanban.go
- publisher/gantt.go
- publisher/todo.go
- publisher/scheduler.go

<a name="license"></a>
## License ##
You can [download DHTMLX components](https://dhtmlx.com/docs/download/) for creating a project management app and test their functionality and compatibility for free during the 30-day trial period. 

However, if you would like to continue using them in your project after the evaluation expires, you should [purchase the license](https://dhtmlx.com/docs/products/licenses.shtml). 

We recommend [exploring the Planning and Complete packs](https://dhtmlx.com/docs/products/licenses.shtml#bundles) that comprise all the widgets showcased in the demo available at a discounted price.

<a name="links"></a>
## Useful links

- [DHTMLX project management tools](https://dhtmlx.com/docs/products/javascript-project-management-library/)
- [More demos with DHTMLX components](https://dhtmlx.com/docs/products/demoApps/)
- [Technical support ](https://forum.dhtmlx.com/c/kanban)
- [Online  documentation](https://docs.dhtmlx.com/)

  <a name="join"></a>
## Join our online community

- Star our GitHub repo :star:
- Keep up with our updates in the [blog](https://dhtmlx.com/blog/) 
- Read us on [Medium](https://dhtmlx.medium.com) :newspaper:
- Follow us on [X](https://x.com/dhtmlx) :bird:
- Check our news and updates on [Facebook](https://www.facebook.com/dhtmlx/) :feet:
