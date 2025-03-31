import {ApolloServer} from '@apollo/server';
import {startStandaloneServer} from '@apollo/server/standalone';
import {ApolloGateway} from '@apollo/gateway';
import {watch} from 'fs';
import {readFile} from 'fs/promises';

const supergraphFilePath = '../graph/supergraph.graphqls';

const server = new ApolloServer({
    gateway: new ApolloGateway({
        async supergraphSdl({update, healthCheck}) {
            const watcher = watch(supergraphFilePath);
            watcher.on('change', async () => {
                try {
                    const updatedSupergraph = await readFile(supergraphFilePath, 'utf-8');
                    await healthCheck(updatedSupergraph);
                    update(updatedSupergraph);
                } catch (e) {
                    console.error(e);
                }
            });
            return {
                supergraphSdl: await readFile(supergraphFilePath, 'utf-8'),
                async cleanup() {
                    watcher.close();
                },
            };
        },
    }),
});

const {url} = await startStandaloneServer(server, {
    listen: {port: 4000},
});

console.log(`🚀  Server ready at: ${url}`);
