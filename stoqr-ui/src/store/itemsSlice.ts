import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { RootState, AppThunk } from './store';
import { Item } from '../model/item';
import axios from 'axios';

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
  axios.get("http://localhost:8080/api/items").then(result => {
    dispatch(set(result.data))
  }).catch(error => {
    console.log(error);
  });
};

export const createItem = (item: Item): AppThunk => dispatch => {
  axios.post("http://localhost:8080/api/items", item).then(result => {
    dispatch(add(result.data));
  }).catch(error => {
    console.log(error);
  });
};

export const deleteItem = (item: Item): AppThunk => dispatch => {
  axios.delete("http://localhost:8080/api/items/" + item.id).then(() => {
    dispatch(remove(item));
  }).catch(error => {
    console.log(error);
  });
};

export const { set, add, remove } = itemsSlice.actions;

export const selectItems = (state: RootState) => state.items;

export default itemsSlice.reducer;
