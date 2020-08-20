import style from './index.module.scss';
import React from 'react';
import { Link } from 'react-router-dom';
import LazyLoad from 'react-lazyload';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

export default function Item(props) {
  const p = props.product;

  let tags = p.tags.slice(0, 4).map((t, i) => {
    return (
      <Link to={`/products/q/best-${t}-avatar-maker`}>{t}{ i<3 && ','} &nbsp;</Link>
    )
  });

  let path = `/products/${p.name.replace(/\W+/mgsi, ' ').replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '')}-${p.short_id}`

  return (
    <div key={p.product_id} className={style.product}>
      <div className={style.wrapper}>
        <Link to={path}>
          <LazyLoad className={style.cover} offset={100}>
            <div className={style.cover} style={{backgroundImage: `url(${p.cover_url})`}} />
          </LazyLoad>
        </Link>
        <div className={style.desc}>
          <Link to={path}>
            <div className={style.name}>{p.name}</div>
          </Link>
          <div className={style.tags}>
            <FontAwesomeIcon className={style.icon} icon={['fas', 'tags']} />
            {tags}
          </div>
        </div>
      </div>
    </div>
  )
}
