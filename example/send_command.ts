
import { MiyuuSsh } from 'miyuu_ssh/mod.ts';

const ssh = MiyuuSsh.open();

const config = ssh.ClientConfig.init('miyuu', 'mecchakawaii!');

const client = ssh.Client.connect(
    config,
    {
        addr_octet: [192, 168, 0, 210]
    }
);

const session = ssh.Session.createSession(client);

session.exec('echo hello world! >> love_love.txt');

// dispose resources.
session.close();
client.close();
config.dispose();

// close miyuu_ssh library.
MiyuuSsh.close();
