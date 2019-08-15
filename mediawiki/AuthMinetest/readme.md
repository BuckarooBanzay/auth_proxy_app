

# Source
* https://git.bananach.space/AuthMinetest.git



```php
$wgGroupPermissions['*']['edit'] = false;
$wgGroupPermissions['*']['createaccount'] = false;
$wgGroupPermissions['*']['autocreateaccount'] = true;


wfLoadExtension( 'AuthMinetest' );

$wgAuthManagerAutoConfig['primaryauth'] = [
	MediaWiki\Auth\MinetestPasswordPrimaryAuthenticationProvider::class => [
		'class' => MediaWiki\Auth\MinetestPasswordPrimaryAuthenticationProvider::class,
		'args' => [
			[
				'minetestUrl' => 'http://your.minetest.url',
			]
		],
		'sort' => 0,
	],
	];
```
