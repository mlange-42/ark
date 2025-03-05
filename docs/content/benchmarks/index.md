+++
title = 'Benchmarks'
type = "docs"
weight = 100
description = "An overview of the runtime cost of typical ECS operations."
+++

This chapter gives an overview of the runtime cost of typical Ark operations.
All time information is per entity.
All components used in the benchmarks have two `int64` fields.
Batch operations are performed in batches of 1000 entities.

Absolute numbers are not really meaningful, as they heavily depend on the hardware.
However, all benchmarks run in the CI in the same job and hence on the same machine, and can be compared.

Benchmark code: {{< repo "/tree/main/benchmark/table" "benchmark/table" >}} in the {{< repo "" "GitHub repository" >}}.

Benchmarks are run automatically in the GitHub CI, and are updated on this page on every merge into the `main` branch.
They always reflect the latest development state of Ark.

{{< html >}}
<br/>
<div class="font-small">
{{< /html >}}

{{% include "/generated/_benchmarks.md" %}}

{{< html >}}
</div>
{{< /html >}}
