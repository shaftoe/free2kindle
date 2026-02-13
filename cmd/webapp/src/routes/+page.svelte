<script lang="ts">
	import { resolve } from '$app/paths';
</script>

<nav>
	<ul>
		<li><a href={resolve('/settings')}>Settings</a></li>
		<li><a href={resolve('/articles')}>My List</a></li>
	</ul>
</nav>
