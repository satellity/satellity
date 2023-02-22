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

export const useGenres = (genre) => {
  return useQuery(
      [genre],
      () => api.gist.genres(genre).then((resp) => resp.data),
      {
        cacheTime: 60 * 1000,
        staleTime: 60 * 1000,
      },
  );
};

export const useFaucet = () => {
  return useQuery(
      ['faucet'],
      () => api.chain.list().then((resp) => resp.data),
      {
        cacheTime: 30 * 60 * 1000,
        staleTime: 60 * 60 * 1000,
      },
  );
};

export const useRatios = () => {
  return useQuery(
      ['ratios'],
      () => api.ratio.index().then((resp) => resp.data),
      {},
  );
};
