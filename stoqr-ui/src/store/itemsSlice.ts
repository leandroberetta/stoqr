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
      state.items.forEach((item, index) => {        
        if (item.id === action.payload.id) {
          state.items.splice(index, 1)
        }          
      });
    },
    withdraw: (state, action: PayloadAction<Item>) => {
      state.items.forEach((item, index) => {        
        if (item.id === action.payload.id) {
          state.items[index].actual -= 1; 
        }          
      });
    },
  },
});

export const fetchItems = (filter: string|null): AppThunk => dispatch => {
  axiosInstance.get("api/items", { params: { filter: filter } }).then((result: AxiosResponse<Item[]>) => {
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

export const withdrawItem = (item: Item): AppThunk => dispatch => {
  axiosInstance.get("api/items/withdraw/" + item.id).then(() => {
    dispatch(withdraw(item));
  }).catch((error: AxiosError) => {
    console.log(error);
  });
};

export const { set, add, remove, withdraw } = itemsSlice.actions;

export const selectItems = (state: RootState) => state.items;

export const selectItem = (id: number) => (state: RootState): Item|null => {
  console.log("selectItem");
  state.items.items.forEach(element => {
    console.log(element);
    if (element.id === id) {
      return element;
    }
  });
  return null;
};

export default itemsSlice.reducer;
