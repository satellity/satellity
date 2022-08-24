import React, {useState} from 'react';
import PropTypes from 'prop-types';
import API from 'api/index.js';
import Button from 'components/button.js';

import style from './index.module.scss';

const New = (props) => {
  const i18n = window.i18n;
  const api = new API();

  const {topicId} = props;
  const [body, setBody] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleChange = (e) => {
    const {value} = e.target;
    setBody(value);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (submitting) {
      return;
    }
    setSubmitting(true);
    api.comment.create(body).then((resp) => {
      if (resp.error) {
        setSubmitting(false);
        return;
      }
      setBody('');
      setSubmitting(false);
    });
  };

  if (!api.user.loggedIn()) {
    return (
      <div className={style.custom}>
        {i18n.t('comment.custom')}
      </div>
    );
  }
  return (
    <div className={style.form}>
      <form onSubmit={handleSubmit}>
        <input type='hidden' name='topic_id' defaultValue={topicId} />
        <div>
          <textarea type='text' name='body' minLength='3' required
            placeholder={i18n.t('comment.form.body')} value={body} onChange={handleChange} />
        </div>
        <div className='action'>
          <Button type='submit' classes='submit' disabled={submitting} text={i18n.t('general.submit')}/>
        </div>
      </form>
    </div>
  );
};

New.propTypes = {
  topicId: PropTypes.string,
};

export default New;
