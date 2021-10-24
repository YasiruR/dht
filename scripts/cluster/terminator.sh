for node in $(cat started_nodes.txt); do
    curl -X POST --retry 3 "${node}:52520/internal/terminate"
done
rm started_nodes.txt