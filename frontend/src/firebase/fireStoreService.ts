import { db } from "@/firebase/firebase";
import {
  collection,
  doc,
  getDoc,
  getDocs,
  addDoc,
  updateDoc,
  setDoc,
  deleteDoc,
  DocumentData,
} from "firebase/firestore";

// Create
export const createDocument = async (
  collectionName: string,
  data: DocumentData
) => {
  const collRef = collection(db, collectionName);
  const docRef = await addDoc(collRef, data);
  return { id: docRef.id, ...data };
};

// Read single
export const readDocument = async (collectionName: string, docId: string) => {
  const docRef = doc(db, collectionName, docId);
  const snapshot = await getDoc(docRef);
  return snapshot.exists() ? { id: snapshot.id, ...snapshot.data() } : null;
};

// Update
export const updateDocument = async (
  collectionName: string,
  docId: string,
  data: Partial<DocumentData>
) => {
  const docRef = doc(db, collectionName, docId);
  await setDoc(docRef, data, { merge: true });
  return { id: docId, ...data };
};

// Delete
export const deleteDocument = async (collectionName: string, docId: string) => {
  const docRef = doc(db, collectionName, docId);
  await deleteDoc(docRef);
  return true;
};
