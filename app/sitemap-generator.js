require("babel-register")({
  presets: ["es2015", "react"]
});

const router = require("./sitemap-routes").default;
const Sitemap = require("react-router-sitemap").default;
const axios = require('axios');

const instance = axios.create({
  baseURL: 'https://routinost.com/api/',
  timeout: 10000,
  headers: {'Content-Type': 'application/json'}
});

async function generateSitemap() {

  const topics = await instance.get('topics').then((resp) => {
    if (resp.error) {
      return;
    }
    return resp.data.data;
  });

  const products = await instance.get('products').then((resp) => {
    if (resp.error) {
      return;
    }
    return resp.data.data;
  });

  let topicsMap = [];
  let productsMap = [];
  let productsQMap = [];

  for (let i=0; i<topics.length; i++) {
    topicsMap.push({id: `${topics[i].short_id}-${topics[i].title.replace(/\W+/mgsi, ' ').replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '')}`});
  }

  let keywords  = new Map();
  for (let i=0; i<products.length; i++) {
    productsMap.push({id: `${products[i].name.replace(/\W+/mgsi, ' ').replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '')}-${products[i].short_id}`});
    for (let j=0;j<products[i].tags.length;j++) {
      keywords.set(products[i].tags[j], `best-${products[i].tags[j]}-avatar-maker`);
    }
  }
  keywords.forEach((v, k, m) => {
    productsQMap.push({id: v});
  });

  const paramsConfig = {
    "/topics/:id": topicsMap,
    "/products/:id": productsMap,
    "/products/q/:id": productsQMap,
  };

  return (
    new Sitemap(router)
    .applyParams(paramsConfig)
    .build("https://routinost.com")
    .save("./public/sitemap.xml")
  );
}

generateSitemap();
