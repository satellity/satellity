import style from './index.module.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import Loading from '../components/loading.js';
import Item from './item.js';

export default class Index extends Component {
  constructor(props) {
    super(props);

    this.api = window.api;
    let q = props.match.params.id || '';
    q = q.replace('best-', '').replace('-avatar-maker', '');
    this.state = {
      q: q,
      products: [],
      loading: true,
    };
  }

  componentDidMount() {
    this.api.product.index(this.state.q).then((resp) => {
      if (resp.error) {
        return;
      }
      this.setState({products: resp.data, loading: false});
    });
  }

  componentDidUpdate(prevProps, prevState) {
    if (this.props.match.params.id !== prevProps.match.params.id) {
      let q = this.props.match.params.id || '';
      q = q.replace('best-', '').replace('-avatar-maker', '');
      this.setState({q: q}, () => {
        this.api.product.index(this.state.q).then((resp) => {
          if (resp.error) {
            return;
          }
          this.setState({products: resp.data, loading: false});
        });
      })
    }
  }

  render() {
    const state = this.state;

    const loadingView = (
      <div className={style.loading}>
        <Loading />
      </div>
    );

    const products = state.products.map((p) => {
      return (
        <Item product={p} />
      )
    });

    let header = (
      <h1 className={style.title}> Collections of <span className={style.keyword}>Person Creator</span> For <span role="img" aria-label="Phone Android iOS">📱</span> or <span role="img" aria-label="Web PC Online">💻</span> </h1>
    );

    if (!!state.q) {
      header = (
        <h1 className={style.title}> Collections of <span className={style.keyword}>Person Creator</span> For <span className={style.keyword}>{state.q}</span>, or <Link to='/products'>Visit <span className={style.keyword}>ALL</span> Avatar Maker</Link></h1>
      );
    }

    return (
      <div>
        { header }
        <div className={style.container}>
          { state.loading && loadingView }
          { !state.loading && products }
        </div>
      </div>
    )
  }
}
