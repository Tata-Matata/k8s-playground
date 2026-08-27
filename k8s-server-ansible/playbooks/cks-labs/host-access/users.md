### Disable interactive login shell for user bob


<details>
<summary>Answer</summary>

```
usermod -s /usr/sbin/nologin bob
```

You can see the result in /etc/passwd:

```
bob:x:1001:1001::/home/bob:/usr/sbin/nologin
```

When someone tries to <code>ssh bob@server</code> - authentication may succeed, but instead of getting a shell, Bob gets something like:

This account is currently not available.

and the connection closes.

It's useful for **service accounts** that need to exist and own files/run processes but should not be used for interactive login. The service can still run as *myservice*, but someone cannot normally log in as that account.

</details>


### Check user bob's shell

<details>
<summary>Answer</summary>

```
grep '^bob:' /etc/passwd
```

OR

```
getent passwd bob
```

</details>

### Set password for user

<details>
<summary>Answer</summary>

```
passwd

```

OR for another user

```
sudo passwd bob
```

</details>

### Lock / unlock user's password for password-based authentication

<details>
<summary>Answer</summary>

```
sudo passwd -l bob

```

#### What actually changes?

The password hash in /etc/shadow gets a special prefix, typically !.


Before:

```
bob:$6$somehash...:...

```

After:

```
bob:!$6$somehash...:...

```

The ! means the existing password hash can no longer successfully authenticate Bob.

#### What can Bob still do?

A locked password doesn't necessarily disable the account completely.

- Log in using Bob's password → no
- SSH using password authentication → no
- Processes can still run as Bob
- Bob's files and permissions remain unchanged
- Authentication using some other mechanism may still be possible, depending on configuration
- if Bob has an **SSH public key** configured, password-locking does not necessarily prevent SSH-key authentication.


#### Unlocking

```
sudo passwd -u bob
```

#### Inspect the status

```
sudo passwd -S bob
```

For example:

```
bob L 08/27/2026 0 99999 7 -1
```

L means locked.

</details>

### Delete user / group

<details>
<summary>Answer</summary>

```
deluser
delgroup

```

</details>

### Create user sam with home dir /opt/sam, shell /bin/bash, UID 2328 and primary group admin

<details>
<summary>Answer</summary>

```
useradd -d /opt/sam -s /bin/bash -u 2328 -g admin sam

```
-G - additional group, not primary one

</details>