import style from './index.module.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import LazyLoad from 'react-lazyload';

export default class Index extends Component {
  constructor(props) {
    super(props);

    this.api = window.api;
    this.state = {
      products: [],
    };
  }

  componentDidMount() {
    this.api.product.index().then((resp) => {
      if (resp.error) {
        return;
      }
      this.setState({products: resp.data});
    });
  }

  render() {
    const state = this.state;

    const products = state.products.map((p) => {
      return (
        <div className={style.product}>
          <Link className={style.wrapper} to={`/products/${p.short_id}-${p.name.replace(/\W+/mgsi, ' ').replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '')}`}>
            <LazyLoad className={style.cover} offset={100}>
              <div className={style.cover} style={{backgroundImage: `url(${p.cover_url})`}} />
            </LazyLoad>
            <div className={style.desc}>
              <div className={style.name}>{p.name}</div>
            </div>
          </Link>
        </div>
      )
    });

    return (
      <div className={style.container}>
        {products}
      </div>
    )
  }
}
