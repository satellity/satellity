import style from './new.module.scss';
import React, {useEffect, useState} from 'react';
import {validate} from 'uuid';
import {Navigate} from 'react-router-dom';
import {useCategory} from 'services';
import API from 'api/index.js';
import Loading from 'components/loading.js';
import Button from 'components/button.js';

const New = (props) => {
  const [api] = useState(new API());
  const [i18n] = useState(window.i18n);
  const [loading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [topic, setTopic] = useState({});

  useEffect(() => {
  }, []);

  const handleChange = (e) => {
    const name = e.target.name;
    const value = e.target.type === 'checkbox' ? (e.target.checked ? 'LINK' : 'POST') : e.target.value;
    topic[name] = value;
    setTopic(topic);
  };

  const handleCategoryClick = (e, value) => {
    e.preventDefault();
    topic['category_id'] = value;
    setTopic(topic);
  };

  const submitForm = () => {
    setSubmitting(true);
    if (validate(topic.topic_id)) {
      api.topic.update(topic.topic_id, topic).then((resp) => {
        if (resp.error) {
          return;
        }
        setSubmitting(false);
      });
      return;
    }
    api.topic.create(topic).then((resp) => {
      if (resp.error) {
        return;
      }
      setSubmitting(false);
    });
  };

  if (!api.me.value()) {
    return (
      <Navigate to="/" replace />
    );
  }

  if (loading) {
    return (
      <div className={style.loading}>
        <Loading class='medium'/>
      </div>
    );
  }

  const {isLoading, data} = useCategory();
  let categories = [];
  if (!isLoading) {
    categories = data.map((c) => {
      return (
        <span key={c.category_id} className={`${style.category} ${c.category_id === topic.category_id ? style.active : ''}`}
          onClick={(e) => handleCategoryClick(e, c.category_id)}>
          {c.alias}
        </span>
      );
    });
  }

  let title = <h1>{i18n.t('topic.title.new')}</h1>;
  if (validate(topic.topic_id)) {
    title = <h1>{i18n.t('topic.title.edit', {name: topic.title})}</h1>;
  }

  const form = (
    <form onSubmit={submitForm}>
      <div className={style.categories}>
        {categories}
      </div>
      <div>
        <input type='text' name='title' pattern='.{3,}' required value={topic.title} autoComplete='off' placeholder={i18n.t('topic.placeholder.title')}
          onChange={handleChange} />
      </div>
      <div className={style.upload}>
        <a href='https://imgur.com/upload' target='_blank' rel='noopener noreferrer'>Does not support upload image, please use imgur first.</a>
      </div>
      {
        topic.topic_type === 'LINK' &&
          <div>
            <textarea name='body' rows='2' value={topic.body} onChange={handleChange} className={style.link}
              placeholder={i18n.t('topic.placeholder.url')} />
          </div>
      }
      <div className={style.submit}>
        <Button type='submit' classes='submit' disabled={submitting} text={i18n.t('general.submit')} />
        {
          !submitting &&
            (topic.topic_id === '' || topic.draft) &&
            <span className={style.draft}>{i18n.t('general.draft')}</span>
        }
      </div>
    </form>
  );

  return (
    <div className='container'>
      <main className='column main'>
        <div className={style.form}>
          {title}
          {form}
        </div>
      </main>
      <aside className='column aside'>
        <div className={style.title}>Rules</div>
        <ul className={style.rules} dangerouslySetInnerHTML={{__html: i18n.t('topic.rules')}}></ul>
      </aside>
    </div>
  );
};

// const TOOLBAR = [
//   {icon: 'heading', action: 'heading', identity: ''},
//   {icon: 'bold', action: 'bold', identity: '**'},
//   {icon: 'italic', action: 'italic', identity: '*'},
//   {icon: 'strikethrough', action: 'strikethrough', identity: '~~'},
//   {icon: 'quote-left', action: 'quote', identity: '> '},
//   {icon: 'list-ol', action: 'ol', identity: '1. '},
//   {icon: 'list-ul', action: 'ul', identity: '* '},
// ];


export default New;
