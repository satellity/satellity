import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import showdown from 'showdown';
import API from '../api/index.js';
import style from './style.css';
import SiteWidget from '../components/site-widget.js';
import CommentList from '../comments/index.js';
import CommentNew from '../comments/new.js';

class TopicShow extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.converter = new showdown.Converter();
    this.state = {
      topic_id: props.match.params.id, title: '', body: '', comments_count: 0, created_at: '', user_id: '', is_author: false,
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
      resp.data.body = this.converter.makeHtml(resp.data.body);
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
  let editAction = '';
  if (state.is_author) {
    editAction = (
      <Link to={`/topics/${state.topic_id}/edit`} className={style.edit}>
        <FontAwesomeIcon icon={['far', 'edit']} />
      </Link>
    )
  }
  let comments =  '';
  if (state.comments_count > 0) {
    comments = <CommentList topicId={state.topic_id} />
  }
  return (
    <div className='container'>
      <main className='section main'>
        <div className={style.content}>
          <header className={style.header}>
            <img src={state.user.avatar_url} className={style.avatar} />
            <h1 className={style.title}>
              {state.title}
              {editAction}
            </h1>
            <div className={style.info}>
              {state.category.name} • {state.user.nickname} • <TimeAgo date={state.created_at} />
            </div>
          </header>
          <article className={`md ${style.body}`} dangerouslySetInnerHTML={{__html: state.body}}>
          </article>
        </div>
        {comments}
        <CommentNew topicId={state.topic_id} />
      </main>
      <aside className='section aside'>
        <SiteWidget />
      </aside>
    </div>
  )
}

export default TopicShow;
