+++
title = 'Benchmarks'
type = "docs"
weight = 999
description = "An overview of the runtime cost of typical ECS operations."
+++

This chapter gives an overview of the runtime cost of typical Ark operations.
All time information is per entity.
All components used in the benchmarks have two `int64` fields.
Batch operations are performed in batches of 1000 entities.

Benchmark code: {{< repo "/tree/main/benchmark/table" "benchmark/table" >}} in the {{< repo "" "GitHub repository" >}}.

Benchmarks are run automatically in the GitHub CI, and are updated on this page on every merge into the `main` branch.
They always reflect the latest development state of Ark.

For comparative benchmarks of different Go ECS implementations, see the [go-ecs-benchmarks](https://github.com/mlange-42/go-ecs-benchmarks) repository.

{{< html >}}
<br/>
<div class="font-small">
{{< /html >}}

{{% include "/generated/_benchmarks.md" %}}

{{< html >}}
</div>
{{< /html >}}
