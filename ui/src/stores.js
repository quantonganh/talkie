import { writable } from "svelte/store";

const { subscribe, set } = writable([]);

export const comments = { subscribe }

export const getCommentsByPost = async () => {
    const response = await fetch('TALKIE_URL' + "/comments" + location.pathname);
    const commentsByPost = await response.json();

    if (response.ok) {
        set(commentsByPost)
    }
};