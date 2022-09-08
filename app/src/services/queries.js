import {useQuery} from '@tanstack/react-query';
import API from 'api';

const api = new API();

export const useCategory = () => {
  return useQuery(
      ['categories'],
      () => api.category.index().then((res) => res.data),
      {
        placeholderData: [],
        cacheTime: Infinity,
        staleTime: Infinity,
      },
  );
};

export const useMe = () => {
  return useQuery(
      ['me'],
      () => api.user.local(),
      {
        cacheTime: Infinity,
        staleTime: Infinity,
      },
  );
};
