# manga_server

This is a very simple manga website built with Go (backend) + React (frontend) + 
MongoDB (database). Suits for beginners in these areas to get started.

For a preview, visit [DragonBlade](http://www.dragonblade.xyz).

## Prerequisite
#### Install NodeJS and NPM
1.  Check instructions [here](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm).
1.  Verify installation succeeded:
```bash
node -v
npm -v
```

#### Install MongoDB and Create Database & Collections 
1.  Install and run MongoDB service on your machine.
1.  Create a new database called `manga_server`.
1.  Create collections `mangas` and `users`.


#### Enable Access Control in MongoDB
1.  Create administrator and an additional user in MongoDB.
1.  Add access control to the DB `manga_server`.
Make it to only allow the user added to access `manga_server`. 

For detailed instructions, please visit MongoDB 
[Tutorial](https://docs.mongodb.com/manual/tutorial/enable-authentication/).

#### Add Sample Manga
1.  In DB Collection "mangas", add two entries:
    ```
    id: 1
    name: "ぼくたちは勉強ができない"
    ChapterNo: 2
    ```
    ```
    id: 2
    name: "ドメスティックな彼女"
    ChapterNo: 2
    ```
    For manipulating database content, try using
[MongoDB Compass](https://www.mongodb.com/products/compass).

1.  Copy folder `/static/samplemanga` to `/static/manga` which contains the
2 sample mangas.

*Note: in similar ways, you can add your favorite manga books.*    
  
## Build Frontend
First, change director to `mangaserver_frontend`:
```bash
cd mangaserver_frontend
```

Next, download node dependencies:
```bash
npm install
```

A `node_modules` folder will appear after the command.

Finally, build the frontend part in development mode:
```bash
npm run build
```
or in production mode:
```bash
npm run buildprod
```
Check the built files in `mangaserver_frontend/dist/`
```bash
ll dist/
```
It should contain `index.html` and `main.js`.
## Build Backend
Now, change dir back to `manga_server` and build the Go server.
```bash
cd ..
go build
```

Most likely it will fail with errors of missing dependencies.
You can fix these with `go get <package_path>` command. If you met difficulties,
read this nice [tutorial](https://www.digitalocean.com/community/tutorials/importing-packages-in-go).

After you fixed all the errors, a binary file `manga_server` will appear in 
the directory. 

## Start the Server
Simply run
```bash
./manga_server --dbUsername=<db_username> --dbPassword=<db_password>
```
where `<db_username>` and `<db_password>` represents the user credentials you set up in
section [Enable Access Control in MongoDB](#enable-access-control-in-mongodb). 

Now you can visit `http://localhost` to see the website.

# WIP Features
1.  Optimizing for mobile browsers.
1.  Preloading next pages.
1.  User login & manga comment section.
1.  Request manga survey. 
1.  Manga feature images gallery.