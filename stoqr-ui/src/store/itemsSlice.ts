import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { RootState, AppThunk } from './store';
import { Item } from '../model/item';
import { axiosInstance } from '../service/service';
import { AxiosResponse, AxiosError } from 'axios';

interface ItemsState {
  items: Item[]
}

const initialState: ItemsState = {
  items: [],
};

export const itemsSlice = createSlice({
  name: 'items',
  initialState,
  reducers: {
    set: (state, action: PayloadAction<Item[]>) => {
      state.items = action.payload;
    },
    add: (state, action: PayloadAction<Item>) => {
      state.items.push(action.payload);
    },
    remove: (state, action: PayloadAction<Item>) => {
      state.items.splice(state.items.indexOf(action.payload), 1)
    },
  },
});

export const fetchItems = (): AppThunk => dispatch => {
  axiosInstance.get("api/items").then((result: AxiosResponse<Item[]>) => {
    dispatch(set(result.data))
  }).catch((error: AxiosError) => {
    console.log(error);
  });
};

export const createItem = (item: Item): AppThunk => dispatch => {
  axiosInstance.post("api/items", item).then((result: AxiosResponse) => {
    dispatch(add(result.data));
  }).catch((error: AxiosError) => {
    console.log(error);
  });
};

export const deleteItem = (item: Item): AppThunk => dispatch => {
  axiosInstance.delete("api/items/" + item.id).then(() => {
    dispatch(remove(item));
  }).catch((error: AxiosError) => {
    console.log(error);
  });
};

export const { set, add, remove } = itemsSlice.actions;

export const selectItems = (state: RootState) => state.items;

export default itemsSlice.reducer;
