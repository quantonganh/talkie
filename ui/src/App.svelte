<svelte:head>
	<script src="https://accounts.google.com/gsi/client" async defer on:load={googleLoaded}></script>
</svelte:head>

<script>
    import {onMount} from 'svelte';
    import {handleCredentialResponse} from './google.js';

    import CreateComment from "./CreateComment.svelte";
    import ListComments from "./ListComments.svelte";

    let googleReady = false;
    let mounted = false;

    onMount(async () => {
        mounted = true;
        if (googleReady) {
			displaySignInButton()
        }
    });

    function googleLoaded() {
        googleReady = true;
        if (mounted) {
			displaySignInButton()
        }
    }

    function displaySignInButton() {
        google.accounts.id.initialize({
            client_id: 'GOOGLE_CLIENT_ID',
            callback: handleCredentialResponse
        });
        google.accounts.id.renderButton(
            document.getElementById("signInWithGoogle"),
            { theme: "outline", size: "large" }  // customization attributes
        );
        google.accounts.id.prompt(); // also display the One Tap dialog
    }
</script>

<div id="signInWithGoogle"></div>

<CreateComment innerHTML="Comment"/>

<ListComments />