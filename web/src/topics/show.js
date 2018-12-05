import './index.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import API from '../api/index.js';
import style from './style.css';
import SiteWidget from '../components/site-widget.js';

class TopicShow extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {
      topic_id: '', title: '', body: '', comments_count: 0, created_at: '', user_id: '', is_author: false,
      user: {user_id: '', nickname: ''},
      category: {category_id: '', name: ''}
    };
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('topic', 'layout');
  }

  componentDidMount() {
    const user = this.api.user.me();
    this.api.topic.show(this.props.match.params.id, (resp) => {
      resp.data.is_author = resp.data.user.user_id === user.user_id;
      this.setState(resp.data);
    });
  }

  render() {
    return (
      <View state={this.state} />
    )
  }
}

const View = ({state}) => {
  var editAction = '';
  if (state.is_author) {
    editAction = (
      <Link to={`/topics/${state.topic_id}/edit`} className={style.edit}>
        <FontAwesomeIcon icon={['far', 'edit']} />
      </Link>
    )
  }
  return (
    <div className='container'>
      <main className='section main'>
        <header className={style.header}>
          <h1 className={style.title}>
            {state.title}
            {editAction}
            <img src={state.user.avatar_url} className={style.avatar} />
          </h1>
          <div className={style.detail}>
            {state.category.name} • {state.user.nickname} • <TimeAgo date={state.created_at} />
          </div>
        </header>
        <article className={style.body}>
          {state.body}
        </article>
      </main>
      <aside className='section aside'>
        <SiteWidget />
      </aside>
    </div>
  )
}

export default TopicShow;
