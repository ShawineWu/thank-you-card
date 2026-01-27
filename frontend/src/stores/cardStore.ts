import { create } from 'zustand';
import type { CardStore, Card, SearchFilters, PaginationInfo } from '../types';

interface CardStoreActions {
  setCards: (cards: Card[]) => void;
  addCard: (card: Card) => void;
  updateCard: (cardId: string, updates: Partial<Card>) => void;
  setCurrentCard: (card: Card | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  setFilters: (filters: Partial<SearchFilters>) => void;
  clearFilters: () => void;
  setPagination: (pagination: PaginationInfo | null) => void;
  reset: () => void;
}

const initialState: CardStore = {
  cards: [],
  currentCard: null,
  isLoading: false,
  error: null,
  filters: {},
  pagination: null,
};

export const useCardStore = create<CardStore & CardStoreActions>((set, get) => ({
  ...initialState,

  setCards: (cards: Card[]) => {
    set({ cards, error: null });
  },

  addCard: (card: Card) => {
    const currentCards = get().cards;
    set({ cards: [card, ...currentCards] });
  },

  updateCard: (cardId: string, updates: Partial<Card>) => {
    const currentCards = get().cards;
    const updatedCards = currentCards.map((card) =>
      card.id === cardId ? { ...card, ...updates } : card
    );
    set({ cards: updatedCards });
  },

  setCurrentCard: (card: Card | null) => {
    set({ currentCard: card });
  },

  setLoading: (loading: boolean) => {
    set({ isLoading: loading });
  },

  setError: (error: string | null) => {
    set({ error });
  },

  setFilters: (newFilters: Partial<SearchFilters>) => {
    const currentFilters = get().filters;
    set({ filters: { ...currentFilters, ...newFilters } });
  },

  clearFilters: () => {
    set({ filters: {} });
  },

  setPagination: (pagination: PaginationInfo | null) => {
    set({ pagination });
  },

  reset: () => {
    set(initialState);
  },
}));