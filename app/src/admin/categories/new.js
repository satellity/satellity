import React, {useState} from 'react';
import API from '../../api/index.js';
import Loading from '../../components/loading.js';

const AdminCategoryNew = () => {
  const api = new API();
  const [category, setCategory] = useState({});

  const handleChange = (e) => {
    const target = e.target;
    const name = target.name;
    console.log('target:', name, target.value, category);
    category[name] = target.value;
    setCategory(category);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    api.category.admin.create(category);
  };

  return (
    <div>
      <h1 className='welcome'>
        Create a new category
      </h1>
      <div className='panel'>
        <form onSubmit={handleSubmit}>
          <div>
            <label name='name'>Name *</label>
            <input type='text' name='name' value={category.name} autoComplete='off' onChange={handleChange}/>
          </div>
          <div>
            <label name='alias'>Alias *</label>
            <input type='text' name='alias' value={category.alias} autoComplete='off' onChange={handleChange}/>
          </div>
          <div>
            <label name='position'>Position</label>
            <input type='number' name='position' value={category.position} autoComplete='off' onChange={handleChange}/>
          </div>
          <div>
            <label name='name'>Description *</label>
            <textarea type='text' name='description' value={category.description} onChange={handleChange} key='description'/>
          </div>
          <div className='action'>
            <button className='btn submit blue' disabled={category.submitting}>
              { category.submitting && <Loading class='small white'/> }
              &nbsp;SUBMIT
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AdminCategoryNew;
