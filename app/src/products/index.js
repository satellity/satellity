import style from './index.module.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import LazyLoad from 'react-lazyload';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import Loading from '../components/loading.js';

export default class Index extends Component {
  constructor(props) {
    super(props);

    this.api = window.api;
    this.state = {
      products: [],
      loading: true,
    };
  }

  componentDidMount() {
    this.api.product.index().then((resp) => {
      if (resp.error) {
        return;
      }
      this.setState({products: resp.data, loading: false});
    });
  }

  render() {
    const state = this.state;

    const loadingView = (
      <div className={style.loading}>
        <Loading />
      </div>
    )

    const products = state.products.map((p) => {
      return (
        <div className={style.product}>
          <Link className={style.wrapper} to={`/products/${p.short_id}-${p.name.replace(/\W+/mgsi, ' ').replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '')}`}>
            <LazyLoad className={style.cover} offset={100}>
              <div className={style.cover} style={{backgroundImage: `url(${p.cover_url})`}} />
            </LazyLoad>
            <div className={style.desc}>
              <div className={style.name}>{p.name}</div>
              <div className={style.tags}>
                <FontAwesomeIcon className={style.icon} icon={['fas', 'tags']} />
                {p.tags.join(', ')}
              </div>
            </div>
          </Link>
        </div>
      )
    });

    return (
      <div>
        <h1 className={style.title}> Collections of <span className={style.keyword}>Person Creator</span> For <span role="img" aria-label="Phone Android iOS">📱</span> or <span role="img" aria-label="Web PC Online">💻</span> </h1>
        <div className={style.container}>
          {state.loading && loadingView}
          {!state.loading && products}
        </div>
      </div>
    )
  }
}
