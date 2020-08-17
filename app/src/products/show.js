import style from './show.module.scss';
import React, {Component} from 'react';
import { Link } from 'react-router-dom';
import showdown from 'showdown';
import showdownHighlight from 'showdown-highlight';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Helmet } from 'react-helmet';
import Config from '../components/config.js';
import SiteWidget from '../home/widget.js';
import Loading from '../components/loading.js';

export default class Show extends Component {
  constructor(props) {
    super(props);

    this.api = window.api;
    this.converter = new showdown.Converter({ extensions: ['header-anchors', showdownHighlight] });

    this.state = {
      product_id: props.match.params.id,
      name: '',
      body: '',
      loading: true,
      tags: [],
    };
  }

  componentDidMount() {
    this.api.product.show(this.state.product_id).then((resp) => {
      if (resp.error) {
        return
      }

      let data = resp.data;
      data.loading = false;
      data.html_body = this.converter.makeHtml(data.body);
      this.setState(data);
    });
  }

  render() {
    const state = this.state;
    const loadingView = (
      <div className={style.loading}>
        <Loading class='medium' />
      </div>
    );

    let start = state.body.indexOf('>') + 1;
    let end = state.body.indexOf(';');
    start = start > 0 ? start : 0;
    end = end > 256 ? 256 : end;
    if (end < 0) end = 256;
    const seoView = (
      <Helmet>
        <meta charSet="utf-9" />
        <title>{`${state.name} 👦 👧 👨 👩 - ${Config.Name}`}</title>
        <meta name='description' content={`🥇 ${state.body.substring(start, end)}`} />
        <link rel="canonical" href={`${Config.Host}/products/${state.short_id}-${state.name.replace(/\W+/mgsi, ' ').replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '')}`} />
      </Helmet>
    );

    let tags = state.tags.map((t) => {
      return (
        <Link to={`/products/q/best-${t}-avatar-maker`}>{t}, &nbsp;</Link>
      )
    });

    const productView = (
      <div className={style.product}>
        <div className={style.cover} style={{backgroundImage: `url(${state.cover_url})`}} />
        <div className={style.content}>
          <h1>
            {state.name}
          </h1>
          <div>
            {state.body !== '' && <article className={`md ${style.body}`} dangerouslySetInnerHTML={{__html: state.html_body}} />}
          </div>
          {
            state.source !== '' &&
            <div>
              Address: <a href={state.source} rel='nofollow noopener noreferrer' target='_blank'>{state.source}</a>
            </div>
          }
          <div className={style.tags}>
            <FontAwesomeIcon className={style.icon} icon={['fas', 'tags']} />
            {tags}
          </div>
        </div>
      </div>
    );

    return (
      <div className='container'>
        {!state.loading && seoView}
        <main className='column main'>
          {state.loading && loadingView}
          {!state.loading && productView}
        </main>
        <aside className='column aside'>
          <SiteWidget />
        </aside>
      </div>
    )
  }
}
