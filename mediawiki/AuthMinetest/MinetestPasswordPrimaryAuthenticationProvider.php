<?php
/**
 * This program is a free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, write to the Free Software Foundation, Inc.,
 * 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.
 * http://www.gnu.org/copyleft/gpl.html
 *
 * @file
 * @ingroup Auth
 */

namespace MediaWiki\Auth;

use User;

/**
 * A primary authentication provider that authenticates the user against a remote Minetest site.
 *
 * @ingroup Auth
 * @since 1.27
 */
class MinetestPasswordPrimaryAuthenticationProvider extends AbstractPrimaryAuthenticationProvider {

	/** @var string The URL of the Minetest site we authenticate against. */
	protected $minetestUrl;

	/** @var array */
    /*	protected $tokens = [];*/

    /**
     * @param array $params Settings
     *  - minetestUrl: The URL of the Minetest site we authenticate against.
     */
    public function __construct( $params = [] ) {

		if ( empty( $params['minetestUrl'] ) ) {
			throw new \InvalidArgumentException( 'The minetestUrl parameter missing in the auth configuration' );
		}

		$this->minetestUrl = $params['minetestUrl'];
    }

	public function beginPrimaryAuthentication( array $reqs ) {
		$req = AuthenticationRequest::getRequestByClass( $reqs, PasswordAuthenticationRequest::class );
		if ( !$req ) {
			return AuthenticationResponse::newAbstain();
		}

		if ( $req->username === null || $req->password === null ) {
			return AuthenticationResponse::newAbstain();
		}

		$username = User::getCanonicalName( $req->username, 'usable' );
		if ( $username === false ) {
			return AuthenticationResponse::newAbstain();
		}

		$token = $this->getMinetestUserToken( $req->username,  $req->password );

		if ( $token === false ) {
			return AuthenticationResponse::newAbstain();

		} else {
			return AuthenticationResponse::newPass( $username );
		}
	}

	/**
	 * Prepares a curl handler to use for querying the Minetest web services.
	 *
	 * @param string $url
	 * @return resource
	 */
	protected function getMinetestCurlClient( $url ) {

		$curl = curl_init( $url );

		curl_setopt_array( $curl, [
			CURLOPT_USERAGENT => 'MWAuthMinetestBot/1.0',
			//CURLOPT_NOBODY => false,
			//CURLOPT_HEADER => true,
			CURLOPT_FOLLOWLOCATION => true,
			CURLOPT_MAXREDIRS => 10,
			CURLOPT_RETURNTRANSFER => true,
			CURLOPT_SSL_VERIFYPEER => 1,
			CURLOPT_SSL_VERIFYHOST => 2,
		]);

		return $curl;
	}

	/**
	 * Attempts to authenticate the user against Minetest. Checks if user is authenticated.
	 *
	 * @param string $username
	 * @param string $password
	 * @return bool False on error, true otherwise
	 */
	protected function getMinetestUserToken( $username,  $password ) {

		$curl = $this->getMinetestCurlClient( $this->minetestUrl.'/api/login' );

		$params = json_encode( array(
			'username' => $username,
			'password' => $password,
		) );


		curl_setopt_array( $curl, [
			CURLOPT_POST => true,
			CURLOPT_POSTFIELDS => $params,
			CURLOPT_HTTPHEADER => array('Content-Type:application/json')
		]);

		$ret = curl_exec( $curl );
		$info = curl_getinfo( $curl );
		$error = curl_error( $curl );
		curl_close( $curl );

		if ( !empty( $error ) ) {
			$this->logger->error( 'AuthMinetest: cURL error: '.$error );
			return false;

		} elseif ( $info['http_code'] != 200 ) {
			$this->logger->error( 'AuthMinetest: cURL error: unexpected HTTP response code '.$info['http_code'] );
			return false;

		}


		$json = json_decode($ret);


		return $json->success;
	}

    /**
     * @param null|\User $user
     * @param AuthenticationResponse $response
     */
    public function postAuthentication( $user, AuthenticationResponse $response ) {
		if ( $response->status !== AuthenticationResponse::PASS ) {
			return;
		}
        return;
	}


	public function testUserCanAuthenticate( $username ) {
		return $this->testUserExists( $username );
	}

	public function testUserExists( $username, $flags = User::READ_NORMAL ) {
		// TODO - there is no easy way to do this without additional web services on the Minetest side.
		return false;
	}

	public function providerAllowsPropertyChange( $property ) {
		return true;
	}

	public function providerAllowsAuthenticationDataChange( AuthenticationRequest $req, $checkData = true) {
		return \StatusValue::newGood( 'ignored' );
	}

	public function providerChangeAuthenticationData( AuthenticationRequest $req ) {
		return;
	}

	public function accountCreationType() {
		return self::TYPE_CREATE;
	}

    public function beginPrimaryAccountCreation( $user, $creator, array $reqs ) {
		throw new \BadMethodCallException( 'This should not get called' );
	}

    public function getAuthenticationRequests( $action, array $options ) {
        switch ( $action ) {
            case AuthManager::ACTION_LOGIN:
                return [ new PasswordAuthenticationRequest() ];
            default:
                return [];
        }
    }
}
